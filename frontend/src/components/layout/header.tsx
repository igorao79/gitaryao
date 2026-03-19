"use client";

import Link from "next/link";
import { useState, useRef, useEffect } from "react";
import { useAuth } from "@/hooks/use-auth";
import { getGithubLoginURL, getGoogleLoginURL } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { GitBranch, Plus, LogOut, Settings, User, ChevronDown } from "lucide-react";

export function Header() {
  const { user, loading, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  return (
    <header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container flex h-14 items-center justify-between px-4 mx-auto">
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center gap-2 font-bold text-lg">
            <GitBranch className="size-5" />
            GitServ
          </Link>
          {user && (
            <nav className="flex items-center gap-4 text-sm">
              <Link href="/dashboard" className="text-muted-foreground hover:text-foreground transition-colors">
                Dashboard
              </Link>
            </nav>
          )}
        </div>

        <div className="flex items-center gap-3">
          {loading ? (
            <div className="size-8 rounded-full bg-muted animate-pulse" />
          ) : user ? (
            <>
              <Link href="/new">
                <Button variant="outline" size="sm">
                  <Plus className="size-4 mr-1" />
                  New
                </Button>
              </Link>

              {/* User menu */}
              <div className="relative" ref={menuRef}>
                <button
                  onClick={() => setMenuOpen(!menuOpen)}
                  className="flex items-center gap-1 rounded-full hover:opacity-80 transition-opacity"
                >
                  {user.avatar_url ? (
                    <img src={user.avatar_url} alt={user.username} className="size-8 rounded-full" />
                  ) : (
                    <div className="flex size-8 items-center justify-center rounded-full bg-primary text-primary-foreground text-xs font-medium">
                      {user.username[0].toUpperCase()}
                    </div>
                  )}
                  <ChevronDown className="size-3 text-muted-foreground" />
                </button>

                {menuOpen && (
                  <div className="absolute right-0 top-full mt-2 w-48 rounded-md border bg-popover p-1 shadow-md z-50">
                    <div className="px-2 py-1.5">
                      <p className="text-sm font-medium">{user.username}</p>
                      <p className="text-xs text-muted-foreground">{user.email}</p>
                    </div>
                    <Separator className="my-1" />
                    <Link
                      href="/dashboard"
                      onClick={() => setMenuOpen(false)}
                      className="flex items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent cursor-pointer"
                    >
                      <User className="size-4" /> Your repos
                    </Link>
                    <Link
                      href="/settings"
                      onClick={() => setMenuOpen(false)}
                      className="flex items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent cursor-pointer"
                    >
                      <Settings className="size-4" /> Settings
                    </Link>
                    <Separator className="my-1" />
                    <button
                      onClick={() => { setMenuOpen(false); logout(); }}
                      className="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent cursor-pointer"
                    >
                      <LogOut className="size-4" /> Sign out
                    </button>
                  </div>
                )}
              </div>
            </>
          ) : (
            <div className="flex items-center gap-2">
              <a href={getGithubLoginURL()}>
                <Button variant="outline" size="sm">Sign in with GitHub</Button>
              </a>
              <a href={getGoogleLoginURL()}>
                <Button variant="outline" size="sm">Sign in with Google</Button>
              </a>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
