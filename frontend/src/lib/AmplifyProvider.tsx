"use client";

import { configureAmplify } from "@/lib/amplifyConfig";

configureAmplify();

export function AmplifyProvider({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}
