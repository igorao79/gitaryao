"use client";

import { useEffect, useState, useMemo } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/layout/header";
import { getRepoBlob } from "@/lib/api";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Button } from "@/components/ui/button";
import { FileText, ArrowLeft, Copy, Download } from "lucide-react";
import { toast } from "sonner";
import hljs from "highlight.js";

const EXT_TO_LANG: Record<string, string> = {
  js: "javascript", jsx: "javascript", ts: "typescript", tsx: "typescript",
  py: "python", rb: "ruby", go: "go", rs: "rust", java: "java",
  c: "c", cpp: "cpp", h: "c", hpp: "cpp", cs: "csharp",
  html: "html", css: "css", scss: "scss", less: "less",
  json: "json", yaml: "yaml", yml: "yaml", toml: "toml",
  xml: "xml", svg: "xml", sql: "sql", sh: "bash", bash: "bash",
  zsh: "bash", fish: "bash", ps1: "powershell",
  md: "markdown", mdx: "markdown", txt: "plaintext",
  dockerfile: "dockerfile", makefile: "makefile",
  env: "bash", gitignore: "bash", editorconfig: "ini",
  php: "php", swift: "swift", kt: "kotlin", scala: "scala",
  lua: "lua", r: "r", dart: "dart", vue: "xml",
};

function getLanguage(filename: string): string {
  const lower = filename.toLowerCase();
  if (lower === "dockerfile" || lower === "makefile") return EXT_TO_LANG[lower] || "plaintext";
  const ext = lower.split(".").pop() || "";
  return EXT_TO_LANG[ext] || "plaintext";
}

function isImageFile(filename: string): boolean {
  const ext = filename.toLowerCase().split(".").pop() || "";
  return ["png", "jpg", "jpeg", "gif", "webp", "svg", "ico", "bmp"].includes(ext);
}

function isBinaryFile(filename: string): boolean {
  const ext = filename.toLowerCase().split(".").pop() || "";
  return ["zip", "gz", "tar", "rar", "7z", "exe", "dll", "so", "dylib",
    "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx",
    "mp3", "mp4", "avi", "mkv", "mov", "wav", "flac",
    "woff", "woff2", "ttf", "eot", "otf"].includes(ext);
}

export default function BlobPage() {
  const params = useParams();
  const owner = params.owner as string;
  const repo = params.repo as string;
  const pathParts = params.path as string[];

  const ref = pathParts[0];
  const filePath = pathParts.slice(1).join("/");
  const fileName = pathParts[pathParts.length - 1];

  const [content, setContent] = useState<string | null>(null);
  const [size, setSize] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isBinaryFile(fileName)) {
      setError("Binary file cannot be displayed");
      setLoading(false);
      return;
    }

    getRepoBlob(owner, repo, ref, filePath)
      .then((data) => {
        setContent(data.content);
        setSize(data.size);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [owner, repo, ref, filePath, fileName]);

  const highlighted = useMemo(() => {
    if (!content) return null;
    const lang = getLanguage(fileName);
    try {
      if (lang !== "plaintext" && hljs.getLanguage(lang)) {
        return hljs.highlight(content, { language: lang }).value;
      }
      return hljs.highlightAuto(content).value;
    } catch {
      return null;
    }
  }, [content, fileName]);

  const lines = content?.split("\n") || [];

  const copyContent = () => {
    if (content) {
      navigator.clipboard.writeText(content);
      toast.success("Copied to clipboard!");
    }
  };

  return (
    <>
      <Header />
      <main className="container mx-auto max-w-5xl px-4 py-8">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 mb-4 text-sm">
          <Link href={`/${owner}/${repo}`} className="text-muted-foreground hover:text-foreground flex items-center gap-1">
            <ArrowLeft className="size-3" />
            {owner}/{repo}
          </Link>
          <span className="text-muted-foreground">/</span>
          <span className="text-foreground font-medium">{filePath}</span>
        </div>

        {loading ? (
          <Card>
            <CardContent className="p-6">
              <Skeleton className="h-4 w-48 mb-4" />
              <Skeleton className="h-64 w-full" />
            </CardContent>
          </Card>
        ) : error ? (
          <Card>
            <CardContent className="py-12 text-center">
              <FileText className="size-12 mx-auto text-muted-foreground mb-4" />
              <p className="text-muted-foreground">{error}</p>
            </CardContent>
          </Card>
        ) : isImageFile(fileName) ? (
          <Card>
            <CardContent className="p-6 flex justify-center">
              <p className="text-sm text-muted-foreground">Image preview not available for git blob content</p>
            </CardContent>
          </Card>
        ) : (
          <Card className="overflow-hidden">
            {/* File header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-neutral-800 bg-neutral-900/50">
              <div className="flex items-center gap-2 text-sm">
                <FileText className="size-4 text-neutral-500" />
                <span className="font-medium">{fileName}</span>
                <span className="text-neutral-500">{formatSize(size)}</span>
                <span className="text-neutral-500">{lines.length} lines</span>
              </div>
              <div className="flex items-center gap-1">
                <Button variant="ghost" size="icon-sm" onClick={copyContent} title="Copy">
                  <Copy className="size-3.5" />
                </Button>
              </div>
            </div>

            {/* Code content */}
            <div className="overflow-x-auto">
              <table className="w-full text-sm font-mono">
                <tbody>
                  {lines.map((line, i) => (
                    <tr key={i} className="hover:bg-neutral-800/30">
                      <td className="select-none text-right text-neutral-600 px-4 py-0 leading-6 w-12 align-top border-r border-neutral-800">
                        {i + 1}
                      </td>
                      <td className="px-4 py-0 leading-6 whitespace-pre">
                        {highlighted ? (
                          <span
                            dangerouslySetInnerHTML={{
                              __html: highlighted.split("\n")[i] || "",
                            }}
                          />
                        ) : (
                          line
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
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
