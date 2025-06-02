import React from "react";

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="flex h-screen overflow-y-auto xl:flex">
      <section className="h-fit min-h-screen w-full grow overflow-auto px-8 py-16 mx-auto max-w-[700px]">
        {children}
      </section>
    </div>
  );
};

export default RootLayout;
