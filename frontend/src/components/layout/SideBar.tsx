"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, FolderKanban, Shield, Key, LogOut } from "lucide-react";
import { signOut } from "aws-amplify/auth";
import { useRouter } from "next/navigation";

const NAV_ITEMS = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/projects", label: "Projects", icon: FolderKanban },
  { href: "/vulnerabilities", label: "Vulnerabilities", icon: Shield },
  { href: "/settings/tokens", label: "API Tokens", icon: Key },
];

export function SideBar() {
  const pathname = usePathname();
  const router = useRouter();

  async function handleLogout() {
    await signOut();
    router.push("/login");
  }

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <Link href="/dashboard" className="sidebar-logo">
          <span> Vulnly</span>
        </Link>
      </div>
      <nav className="sidebar-nav">
        {NAV_ITEMS.map((item) => {
          const Icon = item.icon;
          const isActive =
            pathname === item.href ||
            (item.href !== "/dashboard" && pathname.startsWith(item.href));
          return (
            <Link
              key={item.href}
              href={item.href}
              className={`sidebar-nav-item ${isActive ? "active" : ""}`}
            >
              <Icon />
              <span>{item.label}</span>
            </Link>
          );
        })}
      </nav>
      <div style={{ padding: "1rem", borderTop: "1px solid #333" }}>
        <button
          onClick={handleLogout}
          className="sidebar-nav-item"
          style={{ width: "100%", background: "none", border: "none" }}
        >
          <LogOut />
          <span>Logout</span>
        </button>
      </div>
    </aside>
  );
}
