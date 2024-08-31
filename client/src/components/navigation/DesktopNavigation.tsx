import { memo } from "react";
import { Link } from "react-router-dom";
import { NavLinkProps, NavLinks } from ".";
import LinkButton from "../ui/LinkButton";

interface DesktopNavigationProps {
  menuItems: NavLinks;
}

const DesktopNavigation = memo(({ menuItems }: DesktopNavigationProps) => (
  <nav className="flex items-center gap-8">
    <ul className="space-x-8">
      {menuItems.map((link) => (
        <DesktopNavLink key={link.label} {...link} />
      ))}
    </ul>
    <div className="flex gap-2">
      <LinkButton to="/login" variant="outline" size="sm">
        Log in
      </LinkButton>
      <LinkButton to="/signup" size="sm">
        Sign up
      </LinkButton>
    </div>
  </nav>
));

const DesktopNavLink = ({ href, label }: NavLinkProps) => (
  <Link
    to={href}
    className="text-slate-800 hover:text-purple-500 transition-colors"
  >
    {label}
  </Link>
);

export default DesktopNavigation;
