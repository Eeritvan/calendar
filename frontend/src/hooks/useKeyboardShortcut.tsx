import { useEffect } from "react";

interface UseKeyboardShortcutArgs {
  key: string;
  ctrl: boolean;
  onKeyPressed: () => void;
}

// https://dev.to/barrymichaeldoyle/how-to-build-a-custom-react-hook-to-listen-for-keyboard-events-32b4
export function useKeyboardShortcut({
  key,
  ctrl,
  onKeyPressed
}: UseKeyboardShortcutArgs) {
  useEffect(() => {
    function keyDownHandler(e: globalThis.KeyboardEvent) {
      if (e.key === key && e.ctrlKey === ctrl) {
        e.preventDefault();
        onKeyPressed();
      }
    }

    document.addEventListener("keydown", keyDownHandler);

    return () => {
      document.removeEventListener("keydown", keyDownHandler);
    };
  }, []);
}
