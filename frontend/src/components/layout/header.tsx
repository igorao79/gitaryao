"use client";

import Link from "next/link";
import { useState, useRef, useEffect } from "react";
import { useAuth } from "@/hooks/use-auth";
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
    <header className="border-b border-neutral-800 bg-neutral-950/80 backdrop-blur-xl sticky top-0 z-50">
      <div className="container flex h-14 items-center justify-between px-4 mx-auto">
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center gap-2 font-bold text-lg hover:opacity-80 transition-opacity">
            <div className="flex items-center justify-center size-7 rounded-lg bg-gradient-to-br from-violet-600 to-indigo-600">
              <GitBranch className="size-4 text-white" />
            </div>
            GitaryAO
          </Link>
          {user && (
            <nav className="flex items-center gap-4 text-sm">
              <Link href="/dashboard" className="text-neutral-400 hover:text-white transition-colors">
                Dashboard
              </Link>
            </nav>
          )}
        </div>

        <div className="flex items-center gap-3">
          {loading ? (
            <div className="size-8 rounded-full bg-neutral-800 animate-pulse" />
          ) : user ? (
            <>
              <Link href="/new">
                <Button variant="outline" size="sm">
                  <Plus className="size-4 mr-1" />
                  New
                </Button>
              </Link>

              <div className="relative" ref={menuRef}>
                <button
                  onClick={() => setMenuOpen(!menuOpen)}
                  className="flex items-center gap-1 rounded-full hover:opacity-80 transition-opacity"
                >
                  {user.avatar_url ? (
                    <img src={user.avatar_url} alt={user.username} className="size-8 rounded-full ring-2 ring-neutral-800" />
                  ) : (
                    <div className="flex size-8 items-center justify-center rounded-full bg-gradient-to-br from-violet-600 to-indigo-600 text-white text-xs font-medium">
                      {user.username[0].toUpperCase()}
                    </div>
                  )}
                  <ChevronDown className="size-3 text-neutral-500" />
                </button>

                {menuOpen && (
                  <div className="absolute right-0 top-full mt-2 w-52 rounded-lg border border-neutral-800 bg-neutral-900 p-1.5 shadow-xl z-50">
                    <div className="px-3 py-2">
                      <p className="text-sm font-medium">{user.username}</p>
                      <p className="text-xs text-neutral-500">{user.email}</p>
                    </div>
                    <Separator className="my-1" />
                    <Link
                      href="/dashboard"
                      onClick={() => setMenuOpen(false)}
                      className="flex items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-neutral-800 transition-colors"
                    >
                      <User className="size-4 text-neutral-400" /> Your repos
                    </Link>
                    <Link
                      href="/settings"
                      onClick={() => setMenuOpen(false)}
                      className="flex items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-neutral-800 transition-colors"
                    >
                      <Settings className="size-4 text-neutral-400" /> Settings
                    </Link>
                    <Separator className="my-1" />
                    <button
                      onClick={() => { setMenuOpen(false); logout(); }}
                      className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-neutral-800 transition-colors text-red-400"
                    >
                      <LogOut className="size-4" /> Sign out
                    </button>
                  </div>
                )}
              </div>
            </>
          ) : (
            <Link href="/login">
              <Button size="sm" className="bg-violet-600 hover:bg-violet-700 text-white">
                Sign in
              </Button>
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
