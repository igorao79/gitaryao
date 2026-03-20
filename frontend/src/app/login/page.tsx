"use client";

import { Header } from "@/components/layout/header";
import { useAuth } from "@/hooks/use-auth";
import { getGithubLoginURL, getGoogleLoginURL } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Github, Globe } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function LoginPage() {
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
          <h1 className="text-3xl font-bold tracking-tight">Sign in to GitaryAO</h1>
          <p className="text-muted-foreground">Choose your sign-in method</p>
        </div>

        {!loading && !user && (
          <div className="flex flex-col gap-4 w-full max-w-xs">
            <a href={getGithubLoginURL()}>
              <Button size="lg" className="cursor-pointer gap-3 w-full h-12 text-base">
                <Github className="size-5" />
                Continue with GitHub
              </Button>
            </a>
            <a href={getGoogleLoginURL()}>
              <Button size="lg" variant="outline" className="cursor-pointer gap-3 w-full h-12 text-base">
                <Globe className="size-5" />
                Continue with Google
              </Button>
            </a>
          </div>
        )}
      </main>
    </>
  );
}
