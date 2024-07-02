import WebApp from "@twa-dev/sdk";
import {
  Collection, // "collection1" | "collection2" | [...]
  Metadata    // Type1 | Type2 | [...]
} from "../logic/metadata";
import { plan } from "./content";

type APIChat = {
  ID: number;
  Type: "user" | "group" | "supergroup" | "channel";
  Label: string;
  Username?: string;
  PhotoURL?: string;
  CreatedAt: string;
};

type APIRecord = {
  ID: number;
  Chat: APIChat;
  Collection: string;
  Value: string;
  Metadata: Metadata;
  CreatedAt: string;
};

type APIRequest = {
  Secret: string;
  Limit: number;
  Record: APIRecord;
};

type APIActivityRequest = {
  Secret: string;
  /* [... additional data for activity] */
};

type APIStopActivityRequest = {
  Secret: string;
  Start: string;
  End: string;
};

type Value = string; // "whatever1" | "whatever2"

class _BotBackend {
  async Add(
    collection: Collection,
    value: Value,
    metadata?: Metadata,
  ): Promise<APIRecord> {
    if (metadata === undefined) {
      metadata = null;
    }

    const rsp = await fetch("/api/freestorage", {
      method: "PUT",
      body: JSON.stringify({
        Secret: WebApp.initData,
        Limit: 0,
        Record: {
          Chat: {
            PhotoURL:
              WebApp.initDataUnsafe.user?.photo_url ??
              WebApp.initDataUnsafe.chat?.photo_url,
          },
          Collection: collection,
          Metadata: metadata,
          Value: value,
        },
      } as APIRequest),
    });
    const data = await rsp.json();

    return data as APIRecord;
  }

  async Get(collection: Collection, limit: number): Promise<APIRecord[]> {
    const rsp = await fetch("/api/freestorage", {
      method: "POST",
      body: JSON.stringify({
        Secret: WebApp.initData,
        Limit: limit,
        Record: {
          Collection: collection,
        },
      } as APIRequest),
    });
    const data = await rsp.json();

    return data as APIRecord[];
  }

  async GetLast(collection: Collection): Promise<APIRecord | undefined> {
    return (await this.Get(collection, 1)).at(0);
  }

  async Activity(/* [...] */) {
    const rsp = await fetch("/api/activity", {
      method: "POST",
      body: JSON.stringify({
        Secret: WebApp.initData,
        /* [...] */
      } as APIActivityRequest),
    });
    if (rsp.status != 200) {
      console.error("can't start activity");
    }
  }

  async StopActivity(start: Date, end: Date) {
    const rsp = await fetch("/api/activity/stop", {
      method: "POST",
      body: JSON.stringify({
        Secret: WebApp.initData,
        Start: start.toISOString(),
        End: end.toISOString(),
      } as APIStopActivityRequest),
    });
    if (rsp.status != 200) {
      console.error("can't stop activity");
    }
  }
}

export const BotBackend = new _BotBackend();
