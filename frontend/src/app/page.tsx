"use client";
import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { getCurrentUser } from "aws-amplify/auth";

export default function HomePage() {
  const router = useRouter();
  useEffect(() => {
    getCurrentUser()
      .then(() => {
        // Logged in go to dashboard
        router.replace("/dashboard");
      })
      .catch(() => {
        // not logged in go to login
        router.replace("/login");
      });
  }, [router]);
  return (
    <div
      style={{
        minHeight: "100vh",
        background: "#000",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        color: "#fff",
      }}
    >
      Loading...
    </div>
  );
}
