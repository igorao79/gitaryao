"use client";

import { useState, useEffect, useCallback } from "react";
import { getCurrentUser, removeToken, isAuthenticated, type User } from "@/lib/api";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchUser = useCallback(async () => {
    if (!isAuthenticated()) {
      setUser(null);
      setLoading(false);
      return;
    }

    try {
      const u = await getCurrentUser();
      setUser(u);
    } catch {
      removeToken();
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  const logout = useCallback(() => {
    removeToken();
    setUser(null);
    window.location.href = "/";
  }, []);

  return { user, loading, logout, refetch: fetchUser };
}
