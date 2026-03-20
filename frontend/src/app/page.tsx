"use client";

import { Header } from "@/components/layout/header";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@/components/ui/button";
import { GitBranch } from "lucide-react";
import Link from "next/link";
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
          <h1 className="text-4xl font-bold tracking-tight">GitaryAO</h1>
          <p className="max-w-md text-lg text-muted-foreground">
            Self-hosted Git service. Push, clone, and manage your repositories.
          </p>
        </div>

        {!loading && !user && (
          <Link href="/login">
            <Button size="lg" className="cursor-pointer text-base px-8 py-3 h-12">
              Sign in
            </Button>
          </Link>
        )}
      </main>
    </>
  );
}
