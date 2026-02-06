"use client";

import { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { getCurrentUser, signOut, AuthUser } from "aws-amplify/auth";
type AuthState = {
  user: AuthUser | null;
  isLoading: boolean;
  signOut: () => Promise<void>;
};

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    getCurrentUser()
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setIsLoading(false));
  }, []);
  async function handleSignOut() {
    try {
      await signOut();
      setUser(null);
      window.location.href = "/login";
    } catch (err) {
      console.error("Failed to sign out:", err);
    }
  }

  return (
    <AuthContext.Provider value={{ user, isLoading, signOut: handleSignOut }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
