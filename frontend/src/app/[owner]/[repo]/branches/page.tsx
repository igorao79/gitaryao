"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header";
import { NavTabs } from "@/components/repo/nav-tabs";
import { getRepoBranches, type BranchInfo } from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Code, GitBranch, Clock } from "lucide-react";

export default function BranchesPage() {
  const params = useParams();
  const owner = params.owner as string;
  const repo = params.repo as string;

  const [branches, setBranches] = useState<BranchInfo[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getRepoBranches(owner, repo)
      .then(setBranches)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [owner, repo]);

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
          { label: "Commits", href: `/${owner}/${repo}/commits`, icon: <Clock className="size-4" /> },
          { label: "Branches", href: `/${owner}/${repo}/branches`, icon: <GitBranch className="size-4" /> },
        ]} />

        {loading ? (
          <div className="space-y-2">
            {[1, 2].map((i) => (
              <Card key={i}>
                <CardContent className="py-4"><Skeleton className="h-5 w-32" /></CardContent>
              </Card>
            ))}
          </div>
        ) : branches.length === 0 ? (
          <Card>
            <CardContent className="py-12 text-center text-muted-foreground">
              No branches yet.
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-2">
            {branches.map((branch) => (
              <Card key={branch.name}>
                <CardContent className="flex items-center gap-3 py-4">
                  <GitBranch className="size-5 text-muted-foreground" />
                  <span className="font-mono text-sm">{branch.name}</span>
                  {branch.is_default && (
                    <Badge variant="secondary">default</Badge>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </main>
    </>
  );
}
