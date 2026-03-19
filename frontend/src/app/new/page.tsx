"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";
import { createRepo } from "@/lib/api";
import { Header } from "@/components/layout/header";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { Lock, Globe } from "lucide-react";

export default function NewRepo() {
  const { user } = useAuth();
  const router = useRouter();
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [isPrivate, setIsPrivate] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setSubmitting(true);
    try {
      const repo = await createRepo({
        name: name.trim(),
        description: description.trim(),
        is_private: isPrivate,
      });
      toast.success("Repository created!");
      router.push(`/${repo.owner_name}/${repo.name}`);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : "Failed to create repository";
      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Header />
      <main className="container mx-auto max-w-2xl px-4 py-8">
        <form onSubmit={handleSubmit}>
          <Card>
            <CardHeader>
              <CardTitle>Create a new repository</CardTitle>
              <CardDescription>
                A repository contains all project files, including the revision history.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="owner">Owner</Label>
                <Input id="owner" value={user?.username || ""} disabled />
              </div>

              <div className="space-y-2">
                <Label htmlFor="name">Repository name *</Label>
                <Input
                  id="name"
                  placeholder="my-awesome-project"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  pattern="^[a-zA-Z0-9_\-\.]+$"
                  required
                />
                <p className="text-xs text-muted-foreground">
                  Use letters, numbers, hyphens, dots, and underscores.
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="description">Description (optional)</Label>
                <Input
                  id="description"
                  placeholder="Short description of your project"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                />
              </div>

              <div className="space-y-3">
                <Label>Visibility</Label>
                <div className="flex flex-col gap-2">
                  <label className="flex items-center gap-3 rounded-lg border p-3 cursor-pointer hover:bg-muted/50 transition-colors">
                    <input
                      type="radio"
                      name="visibility"
                      checked={!isPrivate}
                      onChange={() => setIsPrivate(false)}
                      className="size-4"
                    />
                    <Globe className="size-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium text-sm">Public</p>
                      <p className="text-xs text-muted-foreground">Anyone can see this repository.</p>
                    </div>
                  </label>
                  <label className="flex items-center gap-3 rounded-lg border p-3 cursor-pointer hover:bg-muted/50 transition-colors">
                    <input
                      type="radio"
                      name="visibility"
                      checked={isPrivate}
                      onChange={() => setIsPrivate(true)}
                      className="size-4"
                    />
                    <Lock className="size-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium text-sm">Private</p>
                      <p className="text-xs text-muted-foreground">Only you can see this repository.</p>
                    </div>
                  </label>
                </div>
              </div>
            </CardContent>
            <CardFooter>
              <Button type="submit" disabled={!name.trim() || submitting} className="w-full">
                {submitting ? "Creating..." : "Create repository"}
              </Button>
            </CardFooter>
          </Card>
        </form>

        {user && name.trim() && (
          <Card className="mt-4">
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground mb-2">After creating, push your code:</p>
              <pre className="bg-muted rounded-lg p-3 text-xs overflow-x-auto">
{`git remote add origin http://localhost:8080/${user.username}/${name.trim()}.git
git push -u origin master`}
              </pre>
            </CardContent>
          </Card>
        )}
      </main>
    </>
  );
}
