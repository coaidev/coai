import { createContext, useContext, useEffect, useState } from "react";
import { Moon, Sun, Monitor } from "lucide-react";

import { Button } from "./ui/button";
import { getMemory, setMemory } from "@/utils/memory.ts";
import { themeEvent } from "@/events/theme.ts";

const defaultTheme: Theme = "system";

export type Theme = "dark" | "light" | "system";

type ThemeProviderProps = {
  defaultTheme?: Theme;
};

type ThemeProviderState = {
  theme: Theme;
  toggleTheme?: () => void;
};

export function activeTheme(theme: Theme) {
  const root = window.document.documentElement;

  root.classList.remove("light", "dark");

  if (theme === "system") {
    theme = window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  }

  root.classList.add(`${theme}`);

  themeEvent.emit(theme);
}

export function getTheme() {
  return (getMemory("theme") as Theme) || defaultTheme;
}

const initialState: ThemeProviderState = {
  theme: defaultTheme,
  toggleTheme: () => {
    const currentTheme = getMemory("theme");
    let newTheme: Theme;
    const root = window.document.documentElement;


    root.classList.remove("dark", "light","system");

    // dark -> light -> system -> dark
    if (currentTheme === "dark") {
      newTheme = "light";
      root.classList.add(`${newTheme}`);

    } else if (currentTheme === "light") {
      newTheme = "system";

    } else {
      newTheme = "dark";
      root.classList.add(`${newTheme}`);

    }


    activeTheme(newTheme);
    setMemory("theme", newTheme);

  },
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  defaultTheme = "system",
  ...props
}: ThemeProviderProps) {
  const { theme } = useTheme();


  useEffect(() => {

    const savedTheme = getTheme();
    
    activeTheme(savedTheme);
  }, [theme]);

  const value = {
    theme,
    setTheme: (theme: Theme) => {
      setMemory("theme", theme);
    },
    toggleTheme: initialState.toggleTheme,
  };

  return <ThemeProviderContext.Provider {...props} value={value} />;
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");

  return context;
};

export function ModeToggle() {
  const { theme, toggleTheme } = useTheme();
  const [isAnimating, setIsAnimating] = useState(false);
  const [currentIconIndex, setCurrentIconIndex] = useState(0);

  const icons = [
    <Moon key="dark" className="h-[1.2rem] w-[1.2rem]" />,
    <Sun key="light" className="h-[1.2rem] w-[1.2rem]" />,
    <Monitor key="system" className="h-[1.2rem] w-[1.2rem]" />,
  ];
 

  useEffect(() => {
    const savedTheme = getTheme();
    const index = icons.findIndex(icon => icon.key === savedTheme);

    setCurrentIconIndex(index);
  }, [theme]);

  const handleClick = () => {
    if (isAnimating) return; // 防止重复点击
    setIsAnimating(true);


    toggleTheme?.();

    setTimeout(() => {

      const nextIconIndex = (currentIconIndex + 1) % icons.length;

      setCurrentIconIndex(nextIconIndex);
      setIsAnimating(false);
    }, 500);
  };

  return (
    <Button
      variant="outline"
      size="icon"
      onClick={handleClick}
      style={{
        transition: "background-color 0.5s ease",
        position: "relative",
        overflow: "hidden",
      }}
    >
      <div
        style={{
          position: "absolute",
          transform: isAnimating ? "scale(0)" : "scale(1)",
          transition: "transform 0.5s ease, opacity 0.5s ease",
        }}
      >
        {icons[currentIconIndex]}
      </div>

    </Button>
  );
}

export default ModeToggle;
