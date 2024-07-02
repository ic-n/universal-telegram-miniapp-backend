import WebApp from "@twa-dev/sdk";
import { createContext, FC, useEffect, useState, ReactNode } from "react";

export function hslToHex(hsl: string): `#${string}` {
  let [h, s, l] = hsl.split(" ").map((str) => parseFloat(str));
  l /= 100;
  const a = (s * Math.min(l, 1 - l)) / 100;
  const f = (n: number): string => {
    const k = (n + h / 30) % 12;
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1);
    return Math.round(255 * color)
      .toString(16)
      .padStart(2, "0"); // convert to Hex and pad with zeros if necessary
  };

  return `#${f(0)}${f(8)}${f(4)}`;
}
  

type themeData = { dark: boolean; variant: string };
export const ThemeDataContext = createContext({
  theme: { dark: true, variant: "" } as themeData,
  update: (t: themeData) => {
    console.error(`theme update ignorred ${t}`);
  },
});

const localStorageThemeKey = "theme-variant";

export const Themed: FC<{
  children: ReactNode;
}> = ({ children }) => {
  const [theme, setTheme] = useState({
    dark: true,
    variant: "default",
  } as themeData);

  useEffect(() => {
    // theme.dark = WebApp.colorScheme === "dark";
    theme.dark = false;
    theme.variant = localStorage.getItem(localStorageThemeKey) ?? "nice";
    updateTheme(theme);

    setTimeout(() => {
      const ta = document.getElementById("themed-application")?.firstChild;
      if (ta != undefined && ta != null) {
        const v = window
          .getComputedStyle(ta as Element)
          .getPropertyValue("--nextui-background");
        const color = hslToHex(v);
        document
          .querySelectorAll<HTMLElement>(":root")[0]
          .style.setProperty("--real-bg", v);
        WebApp.setHeaderColor(color);
        WebApp.setBackgroundColor(color);
      }
    }, 10);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const updateTheme = (t: themeData) => {
    localStorage.setItem(localStorageThemeKey, t.variant);

    setTheme({
      dark: t.dark,
      variant: t.variant,
    });
  };

  return (
    <ThemeDataContext.Provider
      value={{
        theme: theme,
        update: updateTheme,
      }}
    >
      <main className={`${theme.variant}${theme.dark ? "dark" : "light"}`}>
        <div
          id="themed-application"
          className="text-foreground placeholder-foreground fg-foreground bg-background"
        >
          {children}
        </div>
      </main>
    </ThemeDataContext.Provider>
  );
};
