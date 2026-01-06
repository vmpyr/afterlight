import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Alert } from "@/components/ui/alert"
import { Archive, Plus, AlertCircle, Lock, Shield, Lightbulb } from "lucide-react"

interface Vault {
  id: string
  vault_name: string
  created_at: string
}

function generateSalt(): string {
  const array = new Uint8Array(16);
  window.crypto.getRandomValues(array);
  return Array.from(array)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}

export default function Vaults() {
  const [vaults, setVaults] = useState<Vault[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)

  const [vaultName, setVaultName] = useState("")
  const [hint, setHint] = useState("")
  const [error, setError] = useState("")
  const [creating, setCreating] = useState(false)

  const navigate = useNavigate()

  useEffect(() => {
    fetchVaults()
  }, [])

  const fetchVaults = async () => {
    try {
      const res = await fetch("/api/v1/vaults")
      if (res.ok) {
        const data = await res.json()
        setVaults(data || [])
      } else {
        console.error("Failed to fetch vaults")
      }
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateVault = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setCreating(true)

    // Generate Salt client-side
    const salt = generateSalt()

    try {
      const res = await fetch("/api/v1/vaults", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        // Sending name, salt, and hint. created_at is handled by DB.
        body: JSON.stringify({
          vault_name: vaultName,
          kdf_salt: salt,
          hint: hint
        }),
      })

      if (res.ok) {
        const newVault = await res.json()
        setVaults([newVault, ...vaults]) // Add to list
        setShowCreateModal(false)
        setVaultName("")
        setHint("")
      } else {
        setError("Failed to create vault. Please try again.")
      }
    } catch (err) {
      setError("An error occurred while creating vault")
    } finally {
      setCreating(false)
    }
  }

  const handleVaultClick = (vaultId: string) => {
    // Navigate to the artifacts view of this vault
    navigate(`/vaults/${vaultId}/artifacts`)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-muted-foreground animate-pulse">Loading secure storage...</div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Your Vaults</h1>
          <p className="text-muted-foreground">Manage your secure containers.</p>
        </div>

        {!showCreateModal && (
            <Button onClick={() => setShowCreateModal(true)}>
            <Plus className="mr-2 h-4 w-4" />
            New Vault
            </Button>
        )}
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <span className="ml-2">{error}</span>
        </Alert>
      )}

      {showCreateModal && (
        <Card className="border-primary/50 shadow-md">
          <CardHeader>
            <div className="flex items-center gap-2">
                <Shield className="h-5 w-5 text-primary" />
                <CardTitle>Create Secure Vault</CardTitle>
            </div>
            <CardDescription>
              A vault is a container for your secrets. Encryption keys are generated locally.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreateVault} className="grid gap-4">
              <div className="grid gap-2">
                <Label htmlFor="name">Vault Name</Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="e.g. Financial, Crypto, Letters"
                  required
                  autoFocus
                  value={vaultName}
                  onChange={(e) => setVaultName(e.target.value)}
                />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="hint">
                    Hint <span className="text-muted-foreground font-normal">(Optional)</span>
                </Label>
                <div className="relative">
                    <Lightbulb className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                    id="hint"
                    type="text"
                    className="pl-9"
                    placeholder="e.g. The name of my first pet..."
                    value={hint}
                    onChange={(e) => setHint(e.target.value)}
                    />
                </div>
                <p className="text-xs text-muted-foreground">
                    This hint is visible to your beneficiaries to help them decrypt the vault.
                </p>
              </div>

              <div className="flex gap-2 justify-end pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => {
                    setShowCreateModal(false)
                    setVaultName("")
                    setHint("")
                    setError("")
                  }}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={creating}>
                  {creating ? "Securing..." : "Create Vault"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {vaults.length === 0 && !showCreateModal ? (
        <Card className="bg-muted/10 border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-16">
            <Archive className="h-16 w-16 text-muted-foreground/50 mb-4" />
            <h3 className="text-lg font-medium">No vaults yet</h3>
            <p className="text-muted-foreground text-center mb-6 max-w-sm">
              Create your first vault to start securing your important documents and messages.
            </p>
            <Button variant="secondary" onClick={() => setShowCreateModal(true)}>
                Create Vault
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {vaults.map((vault) => (
            <Card
              key={vault.id}
              className="cursor-pointer hover:border-primary/50 transition-all hover:shadow-md group"
              onClick={() => handleVaultClick(vault.id)}
            >
              <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
                <CardTitle className="text-xl font-semibold tracking-tight truncate pr-4">
                    {vault.vault_name}
                </CardTitle>
                <Lock className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors" />
              </CardHeader>
              <CardContent>
                <div className="text-sm text-muted-foreground mt-2">
                  Created {new Date(vault.created_at).toLocaleDateString("en-GB")}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}