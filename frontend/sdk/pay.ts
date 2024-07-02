import WebApp from "@twa-dev/sdk";

type APIRequest = {
  Secret: string;
};

type Payment = {
  ConfirmationToken: string;
  AmountValue: string;
  AmountCurrency: string;
  Description: string;
};

class _PayBackend {
  async Get(tier: "tier1" | "tier2"): Promise<Payment> {
    const rsp = await fetch("/api/pay/" + tier, {
      method: "POST",
      body: JSON.stringify({
        Secret: WebApp.initData,
      } as APIRequest),
    });
    const data = await rsp.json();

    return data as Payment;
  }
}

export const PayBackend = new _PayBackend();
