import { useMediaQuery } from "usehooks-ts";
import DesktopNavigation from "./DesktopNavigation";
import MobileNavigation from "./MobileNavigation";

export interface NavLinkProps {
  label: string;
  description: string;
  href: string;
  fullWidth?: boolean;
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
      "Take a look at our competative pricing models and see what best suits your needs.",
    fullWidth: true,
  },
];

function Navigation() {
  const isDesktop = useMediaQuery("(min-width: 640px)");
  return isDesktop ? (
    <DesktopNavigation menuItems={links} />
  ) : (
    <MobileNavigation menuItems={links} />
  );
}

export default Navigation;
