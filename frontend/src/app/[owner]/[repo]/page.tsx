"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header";
import { getRepoTree, getRepoBranches, getRepoCommits, getCloneURL, type TreeEntry, type BranchInfo } from "@/lib/api";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { NavTabs } from "@/components/repo/nav-tabs";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
} from "@/components/ui/table";
import { File, Folder, GitBranch, Copy, Clock, Code } from "lucide-react";
import { toast } from "sonner";

export default function RepoPage() {
  const params = useParams();
  const owner = params.owner as string;
  const repo = params.repo as string;

  const [tree, setTree] = useState<TreeEntry[]>([]);
  const [branches, setBranches] = useState<BranchInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [defaultBranch, setDefaultBranch] = useState("master");
  const [commitCount, setCommitCount] = useState<number | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const [branchList] = await Promise.allSettled([
          getRepoBranches(owner, repo),
        ]);

        let branch = "master";
        if (branchList.status === "fulfilled") {
          setBranches(branchList.value);
          const def = branchList.value.find((b) => b.is_default);
          if (def) {
            branch = def.name;
            setDefaultBranch(def.name);
          }
        }

        const [treeData, commits] = await Promise.allSettled([
          getRepoTree(owner, repo, branch),
          getRepoCommits(owner, repo, branch),
        ]);

        if (treeData.status === "fulfilled") setTree(treeData.value);
        else setError(treeData.reason?.message || "Failed to load repository");

        if (commits.status === "fulfilled") setCommitCount(commits.value.length);
      } catch (err: unknown) {
        setError(err instanceof Error ? err.message : "Failed to load repository");
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [owner, repo]);

  const cloneURL = getCloneURL(owner, repo);

  const copyCloneURL = () => {
    navigator.clipboard.writeText(cloneURL);
    toast.success("Copied to clipboard!");
  };

  const sortedTree = [...tree].sort((a, b) => {
    if (a.type === b.type) return a.name.localeCompare(b.name);
    return a.type === "tree" ? -1 : 1;
  });

  const commitBadge = commitCount !== null ? (
    <span className="ml-1 inline-flex items-center justify-center min-w-[20px] h-5 px-1.5 rounded-full bg-neutral-800 text-xs font-medium text-neutral-300">
      {commitCount}
    </span>
  ) : null;

  return (
    <>
      <Header />
      <main className="container mx-auto max-w-4xl px-4 py-8">
        <div className="flex items-start justify-between mb-6">
          <div>
            <h1 className="text-xl font-bold">
              <Link href={`/${owner}`} className="text-muted-foreground hover:text-foreground">
                {owner}
              </Link>
              <span className="text-muted-foreground mx-1">/</span>
              <span>{repo}</span>
            </h1>
          </div>
          <Dialog>
            <DialogTrigger asChild>
              <Button variant="outline" size="sm">
                <Copy className="size-4 mr-2" />
                Clone
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Clone repository</DialogTitle>
              </DialogHeader>
              <div className="flex gap-2">
                <code className="flex-1 bg-muted rounded-md px-3 py-2 text-sm font-mono break-all">
                  {cloneURL}
                </code>
                <Button variant="outline" size="icon" onClick={copyCloneURL}>
                  <Copy className="size-4" />
                </Button>
              </div>
              <pre className="bg-muted rounded-md p-3 text-xs mt-2">
                {`git clone ${cloneURL}`}
              </pre>
            </DialogContent>
          </Dialog>
        </div>

        <NavTabs tabs={[
          { label: "Code", href: `/${owner}/${repo}`, icon: <Code className="size-4" /> },
          { label: "Commits", href: `/${owner}/${repo}/commits`, icon: <Clock className="size-4" />, badge: commitBadge },
          { label: "Branches", href: `/${owner}/${repo}/branches`, icon: <GitBranch className="size-4" /> },
        ]} />

        {loading ? (
          <Card>
            <CardContent className="p-0">
              {[1, 2, 3, 4].map((i) => (
                <div key={i} className="flex items-center gap-3 px-4 py-3 border-b last:border-0">
                  <Skeleton className="size-4" />
                  <Skeleton className="h-4 w-48" />
                </div>
              ))}
            </CardContent>
          </Card>
        ) : error ? (
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-muted-foreground">{error}</p>
              <p className="text-sm text-muted-foreground mt-2">
                This repository might be empty. Push some code first!
              </p>
              <pre className="bg-muted rounded-lg p-3 text-xs mt-4 inline-block text-left">
{`git clone ${cloneURL}
cd ${repo}
echo "# ${repo}" > README.md
git add . && git commit -m "init"
git push origin master`}
              </pre>
            </CardContent>
          </Card>
        ) : (
          <Card>
            <Table>
              <TableBody>
                {sortedTree.map((entry) => (
                  <TableRow key={entry.name} className="hover:bg-muted/50">
                    <TableCell className="w-8">
                      {entry.type === "tree" ? (
                        <Folder className="size-4 text-blue-400" />
                      ) : (
                        <File className="size-4 text-neutral-500" />
                      )}
                    </TableCell>
                    <TableCell>
                      {entry.type === "blob" ? (
                        <Link
                          href={`/${owner}/${repo}/blob/${defaultBranch}/${entry.name}`}
                          className="text-sm hover:text-blue-400 hover:underline transition-colors"
                        >
                          {entry.name}
                        </Link>
                      ) : (
                        <span className="text-sm">{entry.name}</span>
                      )}
                    </TableCell>
                    <TableCell className="text-right text-xs text-muted-foreground">
                      {entry.type === "blob" && entry.size > 0 && formatSize(entry.size)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Card>
        )}
      </main>
    </>
  );
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
