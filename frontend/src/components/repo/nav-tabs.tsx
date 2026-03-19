"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";

interface Tab {
  label: string;
  href: string;
  icon?: React.ReactNode;
  badge?: React.ReactNode;
}

export function NavTabs({ tabs }: { tabs: Tab[] }) {
  const pathname = usePathname();

  return (
    <div className="flex gap-1 border-b mb-6">
      {tabs.map((tab) => {
        const isActive = pathname === tab.href;
        return (
          <Link
            key={tab.href}
            href={tab.href}
            className={cn(
              "flex items-center gap-1.5 px-3 py-2 text-sm font-medium border-b-2 -mb-px transition-colors",
              isActive
                ? "border-primary text-foreground"
                : "border-transparent text-muted-foreground hover:text-foreground hover:border-muted-foreground/25"
            )}
          >
            {tab.icon}
            {tab.label}
            {tab.badge}
          </Link>
        );
      })}
    </div>
  );
}
