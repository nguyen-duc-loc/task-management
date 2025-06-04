"use client";

import { logout } from "@/api/actions/auth";
import { Button } from "@/components/ui/button";
import { IconLogout } from "@tabler/icons-react";
import React from "react";

const LogoutBtn = () => {
  return (
    <Button variant="ghost" onClick={async () => await logout()}>
      <IconLogout />
      Logout
    </Button>
  );
};

export default LogoutBtn;
