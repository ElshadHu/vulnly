"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getCurrentUser } from "aws-amplify/auth";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { SideBar } from "@/components/layout/SideBar";
import "./dashboard.css";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60 * 1000,
      refetchOnWindowFocus: true,
      retry: 1,
    },
  },
});

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getCurrentUser()
      .then(() => setLoading(false))
      .catch(() => router.replace("/login"));
  }, [router]);
  if (loading) {
    return (
      <div className="dashboard-wrapper">
        <div style={{ margin: "auto", color: "#a1a1aa" }}>Loading...</div>
      </div>
    );
  }

  return (
    <QueryClientProvider client={queryClient}>
      <div className="dashboard-wrapper">
        <SideBar />
        <main className="dashboard-main">{children}</main>
      </div>
    </QueryClientProvider>
  );
}
