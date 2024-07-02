import WebApp from "@twa-dev/sdk";

type content = {
  data: {
    id: number;
    attributes: exercise;
  }[];
  meta: {
    pagination: {
      page: number;
      pageSize: number;
      pageCount: number;
      total: number;
    };
  };
};

export type plan = {
  id: number;
  // sensored
  createdAt: "2024-06-28T17:38:55.935Z";
  updatedAt: "2024-06-28T17:38:55.935Z";
  publishedAt: "2024-06-28T17:38:55.935Z";
};

export type exercise = {
  id: number;
  // sensored
  createdAt: "2024-06-28T17:38:55.935Z";
  updatedAt: "2024-06-28T17:38:55.935Z";
  publishedAt: "2024-06-28T17:38:55.935Z";
};

export type muscle = {
  id: number;
  // sensored
  createdAt: "2024-06-28T17:38:55.935Z";
  updatedAt: "2024-06-28T17:38:55.935Z";
  publishedAt: "2024-06-28T17:38:55.935Z";
};

export type task = {
  id: 0;
  // sensored
  createdAt: "2024-06-28T17:38:55.935Z";
  updatedAt: "2024-06-28T17:38:55.935Z";
  publishedAt: "2024-06-28T17:38:55.935Z";
};

export type knowledge = {
  id: 0;
  // sensored
  createdAt: "2024-06-28T17:38:55.935Z";
  updatedAt: "2024-06-28T17:38:55.935Z";
  publishedAt: "2024-06-28T17:38:55.935Z";
};

interface anycontent {
  id: number;
}

class paginated<Type extends anycontent> {
  path: string;
  page: number;
  filters?: unknown;
  constructor(path: string, page: number = 1, filters?: unknown) {
    this.path = path;
    this.page = page;
    this.filters = filters;
  }
  async current(): Promise<Type[]> {
    const cb = await ContentBackend._content(
      this.path,
      this.page,
      this.filters,
    );
    let r: Type[] = [];
    for (const d of cb.data) {
      const v = d.attributes as unknown as Type;
      v.id = d.id;
      r.push(v);
    }
    return r;
  }
  prev(): paginated<Type> {
    return new paginated<Type>(this.path, this.page - 1, this.filters);
  }
  next(): paginated<Type> {
    return new paginated<Type>(this.path, this.page + 1, this.filters);
  }
}

class _ContentBackend {
  async _content(path: string, page: number = 1, filters?: unknown) {
    path += `?populate=*&populate[0]=exercises&populate[exercises][populate]=muscles&pagination[page]=${page}`;
    if (filters != undefined && filters != null) {
      for (const [k, v] of Object.entries(filters)) {
        path += "&" + `filters[${k}]=${v}`;
      }
    }

    const rsp = await fetch("/content/" + path, {
      method: "GET",
      headers: {
        "Telegram-Init-Data": WebApp.initData,
      },
    });

    return (await rsp.json()) as content;
  }

  Knowledge(): paginated<knowledge> {
    return new paginated<knowledge>("knowledges");
  }

  Exercises(): paginated<exercise> {
    return new paginated<exercise>("exercises");
  }

  async Exercise(key: string): Promise<exercise> {
    const result = await this._content(`exercises`, 1, { key: key });
    return result.data[0].attributes as exercise;
  }

  Plan(
    gender: "f" | "m",
    duration: "min40" | "min60" | "min90",
    regularity: "for-2" | "for-3" | "for-4",
    variation: number,
  ): paginated<plan> {
    return new paginated<plan>("plans", 1, {
      gender: gender,
      duration: duration,
      regularity: regularity,
      variation: variation,
    });
  }

  Tasks(week: number): paginated<task> {
    return new paginated<task>("tasks", 1, {
      week: week,
    });
  }
}

export const ContentBackend = new _ContentBackend();
