import { useEffect, useState } from "react"
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import Login from "./pages/Login"
import Dashboard from "./pages/Dashboard"
import Layout from "./components/Layout"
import PrivateRoute from "./components/PrivateRoute"

function App() {
  const [user, setUser] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    checkUser()
  }, [])

  const checkUser = async () => {
    try {
      const res = await fetch("/api/v1/auth/me")
      if (res.ok) {
        setUser(await res.json())
      }
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleLogout = async () => {
    await fetch("/api/v1/auth/logout", { method: "POST" })
    setUser(null)
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            user ? <Navigate to="/" replace /> : <Login onLoginSuccess={setUser} />
          }
        />

        <Route element={<PrivateRoute user={user} isLoading={loading} />}>
          <Route element={<Layout user={user} onLogout={handleLogout} />}>
            <Route path="/" element={<Dashboard user={user} />} />
            {/* Future Routes */}
            <Route path="/vaults" element={<div>Vaults Coming Soon</div>} />
            <Route path="/settings" element={<div>Settings Coming Soon</div>} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App