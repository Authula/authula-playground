"use client";

import { redirect } from "next/navigation";

import { useQuery } from "@tanstack/react-query";
import { GetMeResponse } from "authula";

import { authulaClientBrowser } from "@/lib/authula-client-browser";
import { Spinner } from "@/components/ui/spinner";

export default function DashboardLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const { data, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const response = await authulaClientBrowser.getMe<GetMeResponse>();
        return response;
      } catch (error) {
        console.error(error);
        return null;
      }
    },
  });

  if (isLoading) {
    return (
      <div>
        <Spinner />
      </div>
    );
  }

  if (!data) {
    redirect("/auth/sign-in");
  }

  if (!data.user?.emailVerified) {
    redirect(`/auth/email-verification?email=${data.user.email}`);
  }

  return <>{children}</>;
}
