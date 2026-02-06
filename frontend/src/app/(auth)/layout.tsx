"use client";
import "./auth.css";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getCurrentUser } from "aws-amplify/auth";

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    getCurrentUser()
      .then(() => {
        // User is logged in, redirect to dashboard
        router.replace("/dashboard");
      })
      .catch(() => {
        // Not logged in, continue
        setChecking(false);
      });
  }, [router]);
  if (checking) {
    return (
      <div className="auth-container">
        <div className="auth-card">Loading...</div>
      </div>
    );
  }
  return <>{children}</>;
}
