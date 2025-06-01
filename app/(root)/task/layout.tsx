import React from "react";

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="flex h-screen overflow-y-auto xl:flex">
      <section className="h-fit min-h-screen w-full grow overflow-auto px-8 py-16 xl:px-24 xl:py-16 ">
        {children}
      </section>
    </div>
  );
};

export default RootLayout;
