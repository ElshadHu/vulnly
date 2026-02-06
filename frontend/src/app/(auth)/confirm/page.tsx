"use client";

import { Suspense, useState } from "react";
import { useRouter } from "next/navigation";
import { confirmSignUp, resendSignUpCode } from "aws-amplify/auth";
import Link from "next/link";

function ConfirmForm() {
  const router = useRouter();
  const email =
    typeof window !== "undefined" ? sessionStorage.getItem("pendingConfirmEmail") || "" : "";

  const [code, setCode] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [resending, setResending] = useState(false);

  async function handleSubmit(e: React.BaseSyntheticEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await confirmSignUp({ username: email, confirmationCode: code });
      router.push("/login?confirmed=true");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Confirmation failed");
    } finally {
      setLoading(false);
    }
  }

  async function handleResend() {
    setError("");
    setResending(true);

    try {
      await resendSignUpCode({ username: email });
      setError("Code resent! Check your email.");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to resend code");
    } finally {
      setResending(false);
    }
  }

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">Confirm your email</h1>
        <p style={{ color: "#a1a1aa", marginBottom: "1rem", textAlign: "center" }}>
          We sent a verification code to {email}
        </p>

        {error && <div className="auth-error">{error}</div>}

        <form onSubmit={handleSubmit} className="auth-form">
          <div className="form-group">
            <label htmlFor="code">Verification Code</label>
            <input
              id="code"
              type="text"
              value={code}
              onChange={(e) => setCode(e.target.value)}
              required
              autoComplete="one-time-code"
              disabled={loading}
              placeholder="Enter 6-digit code"
            />
          </div>

          <button type="submit" disabled={loading} className="auth-button">
            {loading ? "Confirming..." : "Confirm"}
          </button>
        </form>

        <p className="auth-footer">
          Did not receive the code?{" "}
          <button
            onClick={handleResend}
            disabled={resending}
            style={{
              background: "none",
              border: "none",
              color: "#fff",
              cursor: "pointer",
              textDecoration: "underline",
            }}
          >
            {resending ? "Sending..." : "Resend"}
          </button>
        </p>

        <p className="auth-footer">
          <Link href="/login">Back to login</Link>
        </p>
      </div>
    </div>
  );
}

export default function ConfirmPage() {
  return (
    <Suspense fallback={<div className="auth-container">Loading...</div>}>
      <ConfirmForm />
    </Suspense>
  );
}
