import { useState, useEffect } from "react";
import { useLocation, useNavigate } from "react-router";

export const meta = () => {
  return [
    { title: "settings" },
    { name: "description", content: "Testing react router!" }
  ];
};

const Settings = () => {
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const location = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    const isSettingsOpen = location.hash === "#settings";
    setIsOpen(isSettingsOpen);
  }, [location.hash]);

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      void navigate(location.pathname + location.search, { replace: true });
    }
  };

  if (location.pathname === "/settings") return null;

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 bg-black/50 flex items-center justify-center"
      onClick={handleOverlayClick}
      role="presentation"
    >
      <div className="bg-white text-black p-6 max-w-md w-full mx-4">
        <h2 className="text-xl font-bold">Settings</h2>
        <div>
          <p>okok</p>
        </div>
      </div>
    </div>
  );
};

export default Settings;
