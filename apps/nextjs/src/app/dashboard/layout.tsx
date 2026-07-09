"use client";

import { redirect } from "next/navigation";

import { authulaClientBrowser } from "@/lib/authula-client-browser";
import { Spinner } from "@/components/ui/spinner";

export default function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const { data, isLoading } = authulaClientBrowser.core.useGetMe({
    query: {
      retry: false,
    },
  });

  if (isLoading) {
    return (
      <div className="grid place-items-center p-4">
        <Spinner />
      </div>
    );
  }

  if (!data) {
    console.log("redirecting to sign in page");
    redirect("/auth/sign-in");
  }

  if (!data.user?.emailVerified) {
    console.log("redirecting to email verification page");
    redirect(`/auth/email-verification?email=${data.user.email}`);
  }

  return <>{children}</>;
}
