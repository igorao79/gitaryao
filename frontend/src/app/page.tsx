"use client";

import { Header } from "@/components/layout/header";
import { useAuth } from "@/hooks/use-auth";
import { Button } from "@igorao79/uivix";
import { GradientText } from "@igorao79/uivix";
import { GridBackground } from "@igorao79/uivix";
import { GitBranch, ArrowRight } from "lucide-react";
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
      <main className="relative flex flex-1 flex-col items-center justify-center gap-10 px-4 overflow-hidden">
        <GridBackground
          variant="dots"
          size={30}
          color="rgba(255,255,255,0.15)"
          followMouse
          maskRadius={250}
          className="!absolute inset-0 -z-10"
        />

        <div className="flex flex-col items-center gap-6 text-center">
          <div className="flex items-center justify-center size-20 rounded-3xl bg-gradient-to-br from-violet-600 to-indigo-600 text-white shadow-lg shadow-violet-500/25">
            <GitBranch className="size-10" />
          </div>

          <GradientText
            as="h1"
            animate
            speed={4}
            colors="from-violet-400 via-indigo-400 to-cyan-400"
            className="text-5xl sm:text-6xl font-extrabold tracking-tight"
          >
            GitaryAO
          </GradientText>

          <p className="max-w-lg text-lg text-neutral-400 leading-relaxed">
            Self-hosted Git service. Push, clone, and manage your repositories
            with a clean, modern interface.
          </p>
        </div>

        {!loading && !user && (
          <Link href="/login">
            <Button
              variant="gradient"
              size="lg"
              color="#7c3aed"
              rightIcon={<ArrowRight className="size-5" />}
            >
              Get started
            </Button>
          </Link>
        )}

        <div className="flex gap-8 mt-4 text-sm text-neutral-500">
          <div className="flex items-center gap-2">
            <div className="size-2 rounded-full bg-green-500" />
            Git Smart HTTP
          </div>
          <div className="flex items-center gap-2">
            <div className="size-2 rounded-full bg-blue-500" />
            OAuth Login
          </div>
          <div className="flex items-center gap-2">
            <div className="size-2 rounded-full bg-purple-500" />
            Auto Backup
          </div>
        </div>
      </main>
    </>
  );
}
