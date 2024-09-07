import { Link } from "react-router-dom";
import Navigation from "./navigation";
import Logo from "./Logo.tsx";

function Header() {
  return (
    <header className="flex justify-between items-center shadow-lg p-4">
      <Link to="/" className="group flex items-center gap-2">
        <Logo size="sm" className="group-hover:grayscale" />
        <span className="text-slate-600 italic font-semibold">Advancely</span>
      </Link>
      <Navigation />
    </header>
  );
}

export default Header;
