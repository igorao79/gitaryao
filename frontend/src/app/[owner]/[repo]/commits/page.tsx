"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header";
import { NavTabs } from "@/components/repo/nav-tabs";
import { getRepoCommits, type CommitInfo } from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Code, GitBranch, Clock, GitCommit } from "lucide-react";

export default function CommitsPage() {
  const params = useParams();
  const owner = params.owner as string;
  const repo = params.repo as string;

  const [commits, setCommits] = useState<CommitInfo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getRepoCommits(owner, repo, "master")
      .then(setCommits)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [owner, repo]);

  const commitBadge = !loading && commits.length > 0 ? (
    <span className="ml-1 inline-flex items-center justify-center min-w-[20px] h-5 px-1.5 rounded-full bg-neutral-800 text-xs font-medium text-neutral-300">
      {commits.length}
    </span>
  ) : null;

  return (
    <>
      <Header />
      <main className="container mx-auto max-w-4xl px-4 py-8">
        <h1 className="text-xl font-bold mb-6">
          <Link href={`/${owner}`} className="text-muted-foreground hover:text-foreground">{owner}</Link>
          <span className="text-muted-foreground mx-1">/</span>
          <Link href={`/${owner}/${repo}`} className="hover:text-foreground">{repo}</Link>
        </h1>

        <NavTabs tabs={[
          { label: "Code", href: `/${owner}/${repo}`, icon: <Code className="size-4" /> },
          { label: "Commits", href: `/${owner}/${repo}/commits`, icon: <Clock className="size-4" />, badge: commitBadge },
          { label: "Branches", href: `/${owner}/${repo}/branches`, icon: <GitBranch className="size-4" /> },
        ]} />

        {loading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <Card key={i}>
                <CardContent className="py-4">
                  <Skeleton className="h-5 w-64 mb-2" />
                  <Skeleton className="h-4 w-48" />
                </CardContent>
              </Card>
            ))}
          </div>
        ) : commits.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center text-muted-foreground">
              No commits yet.
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-2">
            {commits.map((commit) => (
              <Card key={commit.hash} className="hover:border-neutral-700 transition-colors">
                <CardContent className="flex items-start gap-3 py-4">
                  <GitCommit className="size-5 text-muted-foreground mt-0.5 shrink-0" />
                  <div className="min-w-0 flex-1">
                    <p className="font-medium text-sm truncate">{commit.message}</p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {commit.author} committed {new Date(commit.date).toLocaleDateString()}
                    </p>
                  </div>
                  <code className="text-xs bg-neutral-800 rounded px-2 py-1 font-mono shrink-0 text-neutral-300">
                    {commit.hash.substring(0, 7)}
                  </code>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </>
  );
}
