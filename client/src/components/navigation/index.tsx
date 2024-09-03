import { useMediaQuery } from "usehooks-ts";
import DesktopNavigation from "./DesktopNavigation";
import MobileNavigation from "./MobileNavigation";
import {useAuth} from "../../hooks/useAuth.tsx";

type ShowWhen = "authenticated" | "unauthenticated" | "always";

export interface NavLinkProps {
  label: string;
  description: string;
  href: string;
  fullWidth?: boolean;
  showWhen?: ShowWhen;
}

export type NavLinks = NavLinkProps[];

const links: NavLinks = [
  {
    label: "Home",
    href: "/",
    description: "There's no place like home!",
  },
  {
    label: "About",
    href: "/",
    description: "Find out a little of who we are and what we can do for you.",
  },
  {
    label: "Pricing",
    href: "/",
    description:
      "Take a look at our competitive pricing models and see what best suits your needs.",
    fullWidth: true,
    showWhen: "unauthenticated",
  },
  {
    label: "Dashboard",
    href: "/dashboard",
    description: "Your dashboard contains all the things you need at the touch of a button.",
    fullWidth: true,
    showWhen: "authenticated",
  },
];

function Navigation() {
  const { isAuthenticated } = useAuth();
  const isDesktop = useMediaQuery("(min-width: 640px)");

  const filteredLinks = links.filter((link) => {
    switch (link.showWhen) {
      case "always":
      case undefined:
        return true;
      case "authenticated":
        return isAuthenticated;
      case "unauthenticated":
        return !isAuthenticated;
      default:
        return false;
    }
  });

  return isDesktop ? (
    <DesktopNavigation menuItems={filteredLinks} isAuthenticated={isAuthenticated} />
  ) : (
    <MobileNavigation menuItems={filteredLinks} isAuthenticated={isAuthenticated} />
  );
}

export default Navigation;
