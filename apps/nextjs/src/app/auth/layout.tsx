"use client";

import { PropsWithChildren } from "react";
import { redirect, usePathname } from "next/navigation";

import { Spinner } from "@/components/ui/spinner";
import { authulaClientBrowser } from "@/lib/authula-client-browser";

export default function AuthLayout({ children }: PropsWithChildren) {
  const pathname = usePathname();

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

  if (data) {
    if (!data.user.emailVerified) {
      if (pathname === "/auth/email-verification") {
        return <>{children}</>;
      }

      redirect(`/auth/email-verification?email=${data.user.email}`);
    }
    redirect("/dashboard");
  }

  return <>{children}</>;
}
