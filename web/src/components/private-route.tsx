import { Navigate, Outlet } from "react-router-dom"
import { Loader2 } from "lucide-react"

interface PrivateRouteProps {
  user: any
  isLoading: boolean
}

export default function PrivateRoute({ user, isLoading }: PrivateRouteProps) {
  if (isLoading) {
    return (
      <div className="flex h-screen w-full items-center justify-center bg-background">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  // If user is logged in, show the child route (Outlet)
  // If not, redirect to Login
  return user ? <Outlet /> : <Navigate to="/login" replace />
}