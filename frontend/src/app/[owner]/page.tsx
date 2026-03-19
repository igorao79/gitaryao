"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { listUserRepos, type Repository } from "@/lib/api";
import { Header } from "@/components/layout/header";
import { Card, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { GitBranch, Lock, Globe, User } from "lucide-react";

export default function UserProfilePage() {
  const params = useParams();
  const owner = params.owner as string;
  const [repos, setRepos] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listUserRepos(owner)
      .then(setRepos)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [owner]);

  return (
    <>
      <Header />
      <main className="container mx-auto max-w-4xl px-4 py-8">
        <div className="flex items-center gap-3 mb-6">
          <div className="size-12 rounded-full bg-muted flex items-center justify-center">
            <User className="size-6 text-muted-foreground" />
          </div>
          <h1 className="text-2xl font-bold">{owner}</h1>
        </div>

        <h2 className="text-lg font-semibold mb-4">Repositories</h2>

        {loading ? (
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
            <CardHeader className="flex flex-col items-center gap-4 py-12">
              <GitBranch className="size-12 text-muted-foreground" />
              <div className="text-center">
                <h3 className="font-semibold text-lg">No public repositories</h3>
                <p className="text-muted-foreground mt-1">
                  {owner} doesn&apos;t have any public repositories yet.
                </p>
              </div>
            </CardHeader>
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
