"use client";

import { Header } from "@/components/layout/header";
import { useAuth } from "@/hooks/use-auth";
import { getGithubLoginURL, getGoogleLoginURL } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { GitBranch, Github, Globe } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && user) {
      router.push("/dashboard");
    }
  }, [user, loading, router]);

  return (
    <>
      <Header />
      <main className="flex flex-1 flex-col items-center justify-center gap-8 px-4">
        <div className="flex flex-col items-center gap-4 text-center">
          <div className="flex items-center justify-center size-16 rounded-2xl bg-primary text-primary-foreground">
            <GitBranch className="size-8" />
          </div>
          <h1 className="text-4xl font-bold tracking-tight">GitServ</h1>
          <p className="max-w-md text-lg text-muted-foreground">
            Self-hosted Git service. Push, clone, and manage your repositories.
          </p>
        </div>

        {!loading && !user && (
          <div className="flex flex-col gap-3 sm:flex-row">
            <a href={getGithubLoginURL()}>
              <Button size="lg" className="gap-2 w-full sm:w-auto">
                <Github className="size-5" />
                Sign in with GitHub
              </Button>
            </a>
            <a href={getGoogleLoginURL()}>
              <Button size="lg" variant="outline" className="gap-2 w-full sm:w-auto">
                <Globe className="size-5" />
                Sign in with Google
              </Button>
            </a>
          </div>
        )}
      </main>
    </>
  );
}
