"use client";

import { RefreshCw, User } from "lucide-react";
import { useQueryClient, useIsFetching } from "@tanstack/react-query";

type HeaderProps = {
  title: string;
};

export function Header({ title }: HeaderProps) {
  const queryClient = useQueryClient();
  const isFetching = useIsFetching();

  function handleRefresh() {
    queryClient.invalidateQueries();
  }

  return (
    <header className="dashboard-header">
      <h1 className="dashboard-header-title">{title}</h1>

      <div className="dashboard-header-actions">
        <button
          onClick={handleRefresh}
          className="header-button"
          aria-label="Refresh data"
          disabled={isFetching > 0}
          style={{ opacity: isFetching > 0 ? 0.5 : 1 }}
        >
          <RefreshCw size={16} className={isFetching > 0 ? "spin" : ""} />
        </button>

        <button className="header-button" aria-label="User menu">
          <User size={16} />
        </button>
      </div>
    </header>
  );
}
