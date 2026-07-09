"use client";

import { useState, useMemo } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Search, UserCog, LogOut, ShieldAlert, CircleDot } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { authulaClientBrowser } from "@/lib/authula-client-browser";

const ROLE_STYLES: Record<string, string> = {
  Admin: "bg-slate-900 text-white hover:bg-slate-900",
  Editor: "bg-blue-100 text-blue-800 hover:bg-blue-100",
  Support: "bg-violet-100 text-violet-800 hover:bg-violet-100",
  Viewer: "bg-slate-100 text-slate-700 hover:bg-slate-100",
};

function initials(name: string) {
  return name
    .split(" ")
    .map((p) => p[0])
    .slice(0, 2)
    .join("")
    .toUpperCase();
}

export default function ImpersonationPanel() {
  const [impersonatingId, setImpersonatingId] = useState<string | null>(null);
  const [query, setQuery] = useState("");

  const { data, error, isLoading, isError } = useQuery({
    queryKey: ["users"],
    queryFn: async () => {
      return authulaClientBrowser.admin.listUsers();
    },
  });

  const impersonatedUser = useMemo(
    () => data?.users.find((u) => u.id === impersonatingId) ?? null,
    [impersonatingId, data],
  );

  const filteredUsers = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return data?.users ?? [];
    return data?.users.filter(
      (u) =>
        u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q),
    );
  }, [query]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Error: {(error as any).message ?? "Error"} </div>;
  }

  return (
    <div className="mx-auto w-full max-w-4xl space-y-4 p-6">
      {/* Impersonation state banner — deliberately loud so it's never mistaken for the admin's own session */}
      {impersonatedUser && (
        <Alert className="border-amber-300 bg-amber-50 text-amber-900">
          <ShieldAlert className="h-4 w-4 text-amber-600" />
          <AlertTitle className="flex items-center gap-2">
            You're viewing as {impersonatedUser.name}
          </AlertTitle>
          <AlertDescription className="flex items-center justify-between gap-4">
            <span className="text-amber-800">
              Actions you take are performed on this account until you stop.
            </span>
            <Button
              size="sm"
              variant="outline"
              className="shrink-0 border-amber-400 bg-white text-amber-900 hover:bg-amber-100"
              onClick={() => setImpersonatingId(null)}
            >
              <LogOut className="mr-1.5 h-3.5 w-3.5" />
              Stop impersonating
            </Button>
          </AlertDescription>
        </Alert>
      )}

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight text-slate-900">
            Users
          </h1>
          <p className="text-sm text-slate-500">
            {data?.users.length} accounts &middot; impersonate a user to see the
            product from their side.
          </p>
        </div>
        <div className="relative w-64">
          <Search className="pointer-events-none absolute left-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
          <Input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search name, email, role"
            className="pl-8"
          />
        </div>
      </div>

      <div className="overflow-hidden rounded-lg border border-slate-200">
        <Table>
          <TableHeader>
            <TableRow className="bg-slate-50 hover:bg-slate-50">
              <TableHead>User</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {filteredUsers?.map((user) => {
              const isImpersonated = user.id === impersonatingId;
              const disableImpersonate =
                Boolean(impersonatingId) && !isImpersonated;

              return (
                <TableRow
                  key={user.id}
                  className={isImpersonated ? "bg-amber-50/60" : undefined}
                >
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <Avatar className="h-8 w-8">
                        <AvatarFallback className="bg-slate-200 text-xs font-medium text-slate-700">
                          {initials(user.name)}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <div className="text-sm font-medium text-slate-900">
                          {user.name}
                        </div>
                        <div className="text-xs text-slate-500">
                          {user.email}
                        </div>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    {isImpersonated ? (
                      <span className="inline-flex items-center gap-1.5 text-xs font-medium text-amber-700">
                        <CircleDot className="h-3 w-3 fill-amber-500 text-amber-500" />
                        Impersonating
                      </span>
                    ) : (
                      <span className="text-xs text-slate-400">&mdash;</span>
                    )}
                  </TableCell>
                  <TableCell className="text-right">
                    {isImpersonated ? (
                      <Button
                        size="sm"
                        variant="outline"
                        className="border-amber-400 text-amber-900 hover:bg-amber-100"
                        onClick={() => setImpersonatingId(null)}
                      >
                        <LogOut className="mr-1.5 h-3.5 w-3.5" />
                        Stop
                      </Button>
                    ) : (
                      <Button
                        size="sm"
                        variant="outline"
                        disabled={disableImpersonate}
                        onClick={() => setImpersonatingId(user.id)}
                      >
                        <UserCog className="mr-1.5 h-3.5 w-3.5" />
                        Impersonate
                      </Button>
                    )}
                  </TableCell>
                </TableRow>
              );
            })}
            {filteredUsers?.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="py-8 text-center text-sm text-slate-500"
                >
                  No users match "{query}".
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
