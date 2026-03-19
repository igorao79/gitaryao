"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/use-auth";
import { listMyRepos, type Repository } from "@/lib/api";
import { Header } from "@/components/layout/header";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Plus, GitBranch, Lock, Globe } from "lucide-react";

export default function Dashboard() {
  const { user, loading: authLoading } = useAuth();
  const router = useRouter();
  const [repos, setRepos] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/");
      return;
    }
    if (user) {
      listMyRepos()
        .then(setRepos)
        .catch(console.error)
        .finally(() => setLoading(false));
    }
  }, [user, authLoading, router]);

  return (
    <>
      <Header />
      <main className="container mx-auto max-w-4xl px-4 py-8">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Your Repositories</h1>
          <Link href="/new">
            <Button>
              <Plus className="size-4 mr-2" />
              New repository
            </Button>
          </Link>
        </div>

        {loading || authLoading ? (
          <div className="grid gap-4">
            {[1, 2, 3].map((i) => (
              <Card key={i}>
                <CardHeader>
                  <Skeleton className="h-5 w-48" />
                  <Skeleton className="h-4 w-72" />
                </CardHeader>
              </Card>
            ))}
          </div>
        ) : repos.length === 0 ? (
          <Card>
            <CardContent className="flex flex-col items-center gap-4 py-12">
              <GitBranch className="size-12 text-muted-foreground" />
              <div className="text-center">
                <h3 className="font-semibold text-lg">No repositories yet</h3>
                <p className="text-muted-foreground mt-1">
                  Create your first repository to get started.
                </p>
              </div>
              <Link href="/new">
                <Button>
                  <Plus className="size-4 mr-2" />
                  Create repository
                </Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {repos.map((repo) => (
              <Link key={repo.id} href={`/${repo.owner_name}/${repo.name}`}>
                <Card className="hover:border-primary/50 transition-colors cursor-pointer">
                  <CardHeader>
                    <div className="flex items-center gap-2">
                      <CardTitle className="text-base">
                        <span className="text-muted-foreground">{repo.owner_name}</span>
                        <span className="text-muted-foreground mx-1">/</span>
                        <span className="text-foreground">{repo.name}</span>
                      </CardTitle>
                      <Badge variant={repo.is_private ? "secondary" : "outline"}>
                        {repo.is_private ? (
                          <><Lock className="size-3 mr-1" /> Private</>
                        ) : (
                          <><Globe className="size-3 mr-1" /> Public</>
                        )}
                      </Badge>
                    </div>
                    {repo.description && (
                      <CardDescription>{repo.description}</CardDescription>
                    )}
                  </CardHeader>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </main>
    </>
  );
}
