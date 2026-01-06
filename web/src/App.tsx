import { useEffect, useState } from "react"
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom"
import Login from "./pages/Login"
import Register from "./pages/Register"
import Dashboard from "./pages/Dashboard"
import Vaults from "./pages/Vaults"
import VaultDetail from "./pages/VaultDetail"
import Layout from "./components/layout"
import PrivateRoute from "./components/private-route"

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

        <Route
          path="/register"
          element={
            user ? <Navigate to="/" replace /> : <Register onRegisterSuccess={setUser} />
          }
        />

        <Route element={<PrivateRoute user={user} isLoading={loading} />}>
          <Route element={<Layout user={user} onLogout={handleLogout} />}>
            <Route path="/" element={<Dashboard user={user} />} />
            <Route path="/vaults" element={<Vaults />} />
            <Route path="/vaults/:id/artifacts" element={<VaultDetail />} />
            <Route path="/settings" element={<div>Settings Coming Soon</div>} />
          </Route>
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App