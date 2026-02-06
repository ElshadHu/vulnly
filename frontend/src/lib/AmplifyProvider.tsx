"use client";

import { configureAmplify } from "@/lib/amplifyConfig";
import { useEffect } from "react";

export function AmplifyProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    configureAmplify();
  }, []);
  return <>{children}</>;
}
