import { useEffect, useState } from "react"
import { useParams, useNavigate } from "react-router-dom"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Alert } from "@/components/ui/alert"
import { Separator } from "@/components/ui/separator"
import { ArrowLeft, Plus, AlertCircle, Lock } from "lucide-react"

const buf2hex = (buffer: ArrayBuffer) => {
  return [...new Uint8Array(buffer)]
    .map(x => x.toString(16).padStart(2, '0'))
    .join('');
}

const buf2base64 = (buffer: ArrayBuffer) => {
  let binary = '';
  const bytes = new Uint8Array(buffer);
  const len = bytes.byteLength;
  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}

async function encryptData(plaintext: string) {
  const enc = new TextEncoder();
  const encoded = enc.encode(plaintext);

  const iv = window.crypto.getRandomValues(new Uint8Array(12));

  const key = await window.crypto.subtle.generateKey(
    { name: "AES-GCM", length: 256 },
    true,
    ["encrypt", "decrypt"]
  );

  const ciphertext = await window.crypto.subtle.encrypt(
    { name: "AES-GCM", iv: iv },
    key,
    encoded
  );

  return {
    iv: buf2hex(iv.buffer),            // Hex string for DB
    blob: buf2base64(ciphertext),      // Base64 string for Go []byte
  };
}

interface Artifact {
  id: string
  vault_id: string
  message_type: string
  encrypted_blob: string
  iv: string
  created_at: string
}

interface ArtifactList {
  vault_name: string
  hint?: string
  artifacts: Artifact[]
  created_at: string
}

export default function VaultDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [artifactList, setArtifactList] = useState<ArtifactList | null>(null)
  const [loading, setLoading] = useState(true)
  const [showCreateModal, setShowCreateModal] = useState(false)

  // Form State
  const [messageType, setMessageType] = useState("TEXT_MESSAGE")
  const [secretMessage, setSecretMessage] = useState("") // CHANGED: Plain text input
  const [error, setError] = useState("")
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    if (id) {
      fetchArtifactList()
    }
  }, [id])

  const fetchArtifactList = async () => {
    try {
      const res = await fetch(`/api/v1/vaults/${id}/artifacts`)
      if (res.ok) {
        const data = await res.json()
        setArtifactList(data)
      } else if (res.status === 404) {
        setError("Vault not found")
      } else {
        setError("Failed to fetch artifact list")
      }
    } catch (err) {
      setError("An error occurred while fetching artifact list")
    } finally {
      setLoading(false)
    }
  }

  const handleCreateArtifact = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setCreating(true)

    try {
      // 1. Client-Side Encryption
      const { iv, blob } = await encryptData(secretMessage)

      // 2. Send to API
      const res = await fetch(`/api/v1/vaults/${id}/artifacts`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          message_type: messageType,
          encrypted_blob: blob, // Sending Base64 string
          iv: iv,               // Sending Hex string
        }),
      })

      if (res.ok) {
        const newArtifact = await res.json()
        if (artifactList) {
          setArtifactList({
            ...artifactList,
            artifacts: [newArtifact, ...artifactList.artifacts] // Add new to top
          })
        }
        setShowCreateModal(false)
        setMessageType("TEXT_MESSAGE")
        setSecretMessage("")
      } else {
        setError("Failed to create artifact")
      }
    } catch (err) {
      setError("An error occurred while creating artifact")
      console.error(err)
    } finally {
      setCreating(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-muted-foreground">Loading vault details...</div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate("/vaults")}>
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              {artifactList?.vault_name || "Vault"}
            </h1>
            <p className="text-sm text-muted-foreground">
              {artifactList?.artifacts.length || 0} secure items
            </p>
          </div>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="mr-2 h-4 w-4" />
          Add Secret
        </Button>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <span className="ml-2">{error}</span>
        </Alert>
      )}

      {showCreateModal && (
        <Card className="border-primary/50">
          <CardHeader>
            <CardTitle>Add New Secret</CardTitle>
            <CardDescription>
              This message will be encrypted in your browser before being sent.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleCreateArtifact} className="grid gap-4">
              <div className="grid gap-2">
                <Label htmlFor="message_type">Type</Label>
                <select
                  id="message_type"
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-1 text-sm"
                  value={messageType}
                  onChange={(e) => setMessageType(e.target.value)}
                >
                  <option value="TEXT_MESSAGE">Text Message</option>
                  <option value="S3_OBJECT_LINK">File Link</option>
                </select>
              </div>

              {/* INPUT CHANGED: Accepts Plain Text now */}
              <div className="grid gap-2">
                <Label htmlFor="secret_message">Secret Content</Label>
                <Input
                  id="secret_message"
                  type="text"
                  placeholder="e.g., My Bitcoin Seed Phrase..."
                  required
                  autoComplete="off"
                  value={secretMessage}
                  onChange={(e) => setSecretMessage(e.target.value)}
                />
                <p className="text-xs text-muted-foreground">
                  <Lock className="inline h-3 w-3 mr-1" />
                  End-to-end encrypted
                </p>
              </div>

              <div className="flex gap-2 justify-end pt-2">
                <Button
                  type="button"
                  variant="ghost"
                  onClick={() => setShowCreateModal(false)}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={creating}>
                  {creating ? "Encrypting..." : "Encrypt & Save"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      <Separator />

      <div className="grid gap-4">
        {artifactList?.artifacts.map((artifact) => (
          <Card key={artifact.id}>
            <CardHeader className="pb-2">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-medium flex items-center gap-2">
                  <Lock className="h-4 w-4 text-green-600" />
                  {artifact.message_type}
                </CardTitle>
                <span className="text-xs text-muted-foreground">
                  {new Date(artifact.created_at).toLocaleDateString('en-GB')}
                </span>
              </div>
            </CardHeader>
            <CardContent>
              <div className="bg-muted/50 p-3 rounded text-xs font-mono break-all text-muted-foreground">
                {artifact.encrypted_blob.substring(0, 50)}...
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}