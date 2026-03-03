(require 'straight) ; https://github.com/radian-software/straight.el
(require 's)
(require 'cl)

(defconst my-straight-profile-path (cdr (assoc straight-current-profile
                                               straight-profiles)))
(defconst my-straight-repos (concat straight-base-dir "straight/repos/"))
(defconst my-straight-build (concat straight-base-dir "straight/" straight-build-dir "/"))
(defconst my-straight-build-cache-el (concat straight-base-dir "straight/" "build-cache.el"))
(defconst my-straight-modified (concat straight-base-dir "straight/" "modified/"))

(defun my-straight-read-profile ()
  (with-temp-buffer
    (insert-file-contents my-straight-profile-path)
    (read (buffer-string))))

(defun my-straight-list-packages ()
  (cl-loop for profile in (my-straight-read-profile)
           collect (car profile)))

(defun my-straight-list-package-directories ()
  (cl-loop for name in (my-straight-list-packages)
           collect (concat my-straight-repos name)))

(defun my-straight-meta ()
  `(("base_dir" . ,straight-base-dir)
    ("repo_dir" . ,my-straight-repos)
    ("build_dir" . ,my-straight-build)
    ("build_cache" . ,my-straight-build-cache-el)
    ("modified_dir" . ,my-straight-modified)))

(defun my-straight-info ()
  `(("meta" . ,(my-straight-meta))
    ("profile" . ,(my-straight-read-profile))
    ("packages" . ,(my-straight-list-packages))
    ("directories" . ,(my-straight-list-package-directories))))

(defun my-straight-write (text)
  (append-to-file text nil "/dev/stdout"))

(defun my-straight-build-all ()
  (straight-rebuild-all))

(defun my-straight-update-packages (pkgs)
  (cl-loop for pkg in pkgs
           do (straight-pull-package-and-deps pkg))
  (my-straight-build-all))

(defun my-straight-update-all ()
  (straight-pull-all)
  (my-straight-build-all))

;;
;; Functions to be called by go
;;

(defun my-straight-write-info ()
  (my-straight-write (json-encode (my-straight-info))))

(defun my-straight-commit ()
  (straight-freeze-versions))

(defun my-straight-rollback ()
  (straight-thaw-versions))

(defun my-straight-update (pkgs)
  (if pkgs (my-straight-update-packages pkgs)
    (straight-pull-all)))
