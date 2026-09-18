-- add a directory to the file watcher
merodi.watch.add("src")

-- called when the file watcher starts
merodi.hook("on_start_watching", function()
	merodi.log.info("Start watching")
end)

-- called when a file system event is detected
merodi.hook("on_file_changed", function(mode, filepath)
	-- only log file write events
	if mode == "WRITE" then
		merodi.log.info("MODE: " .. mode .. ", FILEPATH: " .. filepath)
	end
end)
