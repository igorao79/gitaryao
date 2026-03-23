"use client";

import { Header } from "@/components/layout/header";
import { useAuth } from "@/hooks/use-auth";
import { getGithubLoginURL, getGoogleLoginURL } from "@/lib/api";
import { Card, CardContent } from "@igorao79/uivix";
import { Button } from "@igorao79/uivix";
import { GridBackground } from "@igorao79/uivix";
import { Github, Globe, ArrowLeft } from "lucide-react";
import Link from "next/link";
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
      <main className="relative flex flex-1 flex-col items-center justify-center gap-8 px-4 overflow-hidden">
        <GridBackground
          variant="dots"
          size={30}
          color="rgba(255,255,255,0.08)"
          className="!absolute inset-0 -z-10"
        />

        <Card variant="glass" className="w-full max-w-sm p-0">
          <CardContent className="flex flex-col gap-6 p-8">
            <div className="text-center">
              <h1 className="text-2xl font-bold tracking-tight">Welcome back</h1>
              <p className="text-sm text-neutral-400 mt-2">
                Sign in to continue to GitaryAO
              </p>
            </div>

            {!loading && !user && (
              <div className="flex flex-col gap-3">
                <a href={getGithubLoginURL()}>
                  <Button
                    variant="default"
                    size="lg"
                    leftIcon={<Github className="size-5" />}
                    className="w-full !bg-white !text-black hover:!bg-neutral-200 transition-colors"
                  >
                    Continue with GitHub
                  </Button>
                </a>
                <a href={getGoogleLoginURL()}>
                  <Button
                    variant="outline"
                    size="lg"
                    leftIcon={<Globe className="size-5" />}
                    className="w-full"
                  >
                    Continue with Google
                  </Button>
                </a>
              </div>
            )}

            <div className="text-center">
              <Link
                href="/"
                className="inline-flex items-center gap-1 text-sm text-neutral-400 hover:text-white transition-colors"
              >
                <ArrowLeft className="size-3" />
                Back to home
              </Link>
            </div>
          </CardContent>
        </Card>
      </main>
    </>
  );
}
