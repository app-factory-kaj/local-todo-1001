import type { JSX } from "react";
import { Outlet } from "react-router-dom";
import { AppShell, Header } from "@wso2/oxygen-ui";

// The app shell, per wireframes.dsl's TodoList screen: a navbar branded
// "Todo" and nothing else — no sidebar (one screen, nowhere else to go), no
// user menu or sign-out (no auth dependency, no sign-in).
export default function AppLayout(): JSX.Element {
  return (
    <AppShell>
      <AppShell.Navbar>
        <Header>
          <Header.Brand>
            <Header.BrandTitle>Todo</Header.BrandTitle>
          </Header.Brand>
        </Header>
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}
