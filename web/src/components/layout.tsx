import { Outlet, Link, useLocation, useNavigate } from "react-router-dom"
import { LayoutDashboard, Lock, Settings, LogOut, Flashlight } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { ModeToggle } from "@/components/mode-toggle"

interface LayoutProps {
  user: any
  onLogout: () => void
}

export default function Layout({ user, onLogout }: LayoutProps) {
  const location = useLocation()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await onLogout()
    navigate("/login")
  }

  const navItems = [
    { href: "/", label: "Dashboard", icon: LayoutDashboard },
    { href: "/vaults", label: "Vaults", icon: Lock },
    { href: "/settings", label: "Settings", icon: Settings },
  ]

  return (
    <div className="flex min-h-screen w-full bg-muted/40">
      {/* Sidebar */}
      <aside className="fixed inset-y-0 left-0 z-10 hidden w-64 flex-col border-r bg-background sm:flex">
        <div className="flex h-14 items-center justify-between border-b px-4 font-semibold tracking-wide">
          <div className="flex items-center">
            <Flashlight className="mr-2 h-5 w-5" />
            Afterlight
          </div>
          <ModeToggle />
        </div>
        <div className="flex-1 py-4">
          <nav className="grid items-start px-4 text-sm font-medium">
            {navItems.map((item) => {
              const Icon = item.icon
              const isActive = location.pathname === item.href
              return (
                <Link
                  key={item.href}
                  to={item.href}
                  className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-primary ${
                    isActive
                      ? "bg-muted text-primary"
                      : "text-muted-foreground"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </Link>
              )
            })}
          </nav>
        </div>
        <div className="mt-auto p-4">
          <Separator className="my-4" />
          <div className="flex items-center gap-3 px-2 mb-4">
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-primary-foreground font-bold text-xs">
              {user.name.charAt(0).toUpperCase()}
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-medium">{user.name}</span>
              <span className="text-xs text-muted-foreground truncate w-32">
                {user.email}
              </span>
            </div>
          </div>
          <Button
            variant="outline"
            className="w-full justify-start gap-2 text-destructive hover:text-destructive hover:bg-destructive/10"
            onClick={handleLogout}
          >
            <LogOut className="h-4 w-4" />
            Sign Out
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex flex-1 flex-col sm:ml-64">
        <div className="p-8">
          <Outlet />
        </div>
      </main>
    </div>
  )
}