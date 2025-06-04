import ROUTES from "@/constants/routes";
import { IconStack2 } from "@tabler/icons-react";
import Link from "next/link";
import React from "react";
import ModeToggle from "./_components/ModeToggle";
import LogoutBtn from "./_components/LogoutBtn";

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="h-screen overflow-y-auto px-8">
      <div className="mx-auto max-w-5xl mt-10 flex gap-4 justify-between max-sm:flex-col">
        <Link
          className="font-semibold text-2xl text-primary flex gap-2 items-center"
          href={ROUTES.dashboard}
        >
          <IconStack2 className="size-10" />
          Task Management
        </Link>
        <div className="flex gap-2">
          <ModeToggle />
          <LogoutBtn />
        </div>
      </div>
      <section className="h-fit min-h-screen w-full grow overflow-auto py-12 flex justify-center">
        {children}
      </section>
    </div>
  );
};

export default RootLayout;
