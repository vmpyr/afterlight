import { Moon, Sun, Monitor } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useTheme } from "@/components/theme-provider"

export function ModeToggle() {
  const { theme, setTheme } = useTheme()

  const toggleTheme = () => {
    if (theme === "light") {
      setTheme("dark")
    } else if (theme === "dark") {
      setTheme("system")
    } else {
      setTheme("light")
    }
  }

  const getIcon = () => {
    return (
      <>
        <Sun className={`h-4 w-4 transition-all duration-300 ${
          theme === "light" 
            ? "rotate-0 scale-100" 
            : "rotate-90 scale-0"
        }`} />
        <Moon className={`absolute h-4 w-4 transition-all duration-300 ${
          theme === "dark" 
            ? "rotate-0 scale-100" 
            : "-rotate-90 scale-0"
        }`} />
        <Monitor className={`absolute h-4 w-4 transition-all duration-300 ${
          theme === "system" 
            ? "rotate-0 scale-100" 
            : "rotate-180 scale-0"
        }`} />
      </>
    )
  }

  const getTooltip = () => {
    if (theme === "light") return "Switch to dark mode"
    if (theme === "dark") return "Switch to system mode"
    return "Switch to light mode"
  }

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={toggleTheme}
      className="h-9 w-9"
      title={getTooltip()}
    >
      {getIcon()}
      <span className="sr-only">Toggle theme</span>
    </Button>
  )
}
