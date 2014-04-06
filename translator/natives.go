package translator

var pkgNatives = make(map[string]map[string]string)

func init() {
	pkgNatives["os"] = map[string]string{
		"toplevel": `
			if (go$packages["syscall"].Syscall15 !== undefined) { // windows
				NewFile = go$pkg.NewFile = function() { return new File.Ptr(); };
			}
		`,
		"toplevelDependencies": `syscall:Syscall15 os:NewFile`,
		"Args":                 `new (go$sliceType(Go$String))((typeof process !== 'undefined') ? process.argv.slice(1) : [])`,
	}

	pkgNatives["runtime"] = map[string]string{
		"toplevel": `
			go$throwRuntimeError = function(msg) { throw go$panic(new errorString(msg)); };
		`,
		"toplevelDependencies": `runtime:errorString runtime:TypeAssertionError`,
		"getgoroot": `function() {
			return (typeof process !== 'undefined') ? (process.env["GOROOT"] || "") : "/";
		}`,
		"sizeof_C_MStats": `3712`,
		"Caller": `function(skip) {
			var line = go$getStack()[skip + 3];
			if (line === undefined) {
				return [0, "", 0, false];
			}
			var parts = line.substring(line.indexOf("(") + 1, line.indexOf(")")).split(":");
			return [0, parts[0], parseInt(parts[1]), true];
		}`,
		"GC": `function() {}`,
		"GOMAXPROCS": `function(n) {
			if (n > 1) {
				go$notSupported("GOMAXPROCS != 1");
			}
			return 1;
		}`,
		"Goexit": `function() {
			var err = new Go$Error();
			err.go$exit = true;
			throw err;
		}`,
		"NumCPU":       `function() { return 1; }`,
		"ReadMemStats": `function() {}`,
		"SetFinalizer": `function() {}`,
	}

	pkgNatives["sync"] = map[string]string{
		"copyChecker.check":    `function() {}`,
		"runtime_Syncsemcheck": `function() {}`,
	}

	pkgNatives["syscall"] = map[string]string{
		"toplevel": `
			if (go$pkg.Syscall15 !== undefined) { // windows
				Syscall = Syscall6 = Syscall9 = Syscall12 = Syscall15 = go$pkg.Syscall = go$pkg.Syscall6 = go$pkg.Syscall9 = go$pkg.Syscall12 = go$pkg.Syscall15 = loadlibrary = getprocaddress = function() { throw new Error("Syscalls not available."); };
				getStdHandle = GetCommandLine = go$pkg.GetCommandLine = function() {};
				CommandLineToArgv = go$pkg.CommandLineToArgv = function() { return [null, {}]; };
				Getenv = go$pkg.Getenv = function(key) { return ["", false]; };
				GetTimeZoneInformation = go$pkg.GetTimeZoneInformation = function() { return [undefined, true]; };
			} else if (typeof process === "undefined") {
				var syscall = function() { throw new Error("Syscalls not available."); };
				if (typeof go$syscall !== "undefined") {
					syscall = go$syscall;
				}
				Syscall = Syscall6 = RawSyscall = RawSyscall6 = go$pkg.Syscall = go$pkg.Syscall6 = go$pkg.RawSyscall = go$pkg.RawSyscall6 = syscall;
				envs = new (go$sliceType(Go$String))(new Array(0));
			} else {
				try {
					var syscall = require("syscall");
					Syscall = go$pkg.Syscall = syscall.Syscall;
					Syscall6 = go$pkg.Syscall6 = syscall.Syscall6;
					RawSyscall = go$pkg.RawSyscall = syscall.Syscall;
					RawSyscall6 = go$pkg.RawSyscall6 = syscall.Syscall6;
				} catch (e) {
					Syscall = Syscall6 = RawSyscall = RawSyscall6 = go$pkg.Syscall = go$pkg.Syscall6 = go$pkg.RawSyscall = go$pkg.RawSyscall6 = function() { throw e; };
				}
				BytePtrFromString = go$pkg.BytePtrFromString = function(s) { return [go$stringToBytes(s, true), null]; };

				var envkeys = Object.keys(process.env);
				envs = new (go$sliceType(Go$String))(new Array(envkeys.length));
				var i;
				for(i = 0; i < envkeys.length; i++) {
					envs.array[i] = envkeys[i] + "=" + process.env[envkeys[i]];
				}
			}
		`,
		"toplevelDependencies": `syscall:Syscall syscall:Syscall6 syscall:Syscall9 syscall:Syscall12 syscall:Syscall15 syscall:loadlibrary syscall:getprocaddress syscall:getStdHandle syscall:GetCommandLine syscall:CommandLineToArgv syscall:Getenv syscall:GetTimeZoneInformation syscall:RawSyscall syscall:RawSyscall6 syscall:BytePtrFromString`,
	}

	pkgNatives["testing"] = map[string]string{
		"runTest": `function(f, t) {
			try {
				f(t);
				return null;
			} catch (e) {
				go$jsErr = null;
				return e;
			}
		}`,
	}

	pkgNatives["time"] = map[string]string{
		"now":   `go$now`,
		"After": `function() { go$notSupported("time.After (use time.AfterFunc instead)") }`,
		"AfterFunc": `function(d, f) {
			setTimeout(f, go$div64(d, new Duration(0, 1000000)).low);
			return null;
		}`,
		"NewTimer": `function() { go$notSupported("time.NewTimer (use time.AfterFunc instead)") }`,
		"Sleep":    `function() { go$notSupported("time.Sleep (use time.AfterFunc instead)") }`,
		"Tick":     `function() { go$notSupported("time.Tick (use time.AfterFunc instead)") }`,
	}

	pkgNatives["github.com/gopherjs/gopherjs/js"] = map[string]string{
		"toplevelDependencies": `github.com/gopherjs/gopherjs/js:Error`,
	}

	pkgNatives["github.com/gopherjs/gopherjs/js_test"] = map[string]string{
		"dummys": `{
			someBool: true,
			someString: "abc\u1234",
			someInt: 42,
			someFloat: 42.123,
			someArray: [41, 42, 43],
			add: function(a, b) {
				return a + b;
			},
			mapArray: go$mapArray,
			toUnixTimestamp: function(d) {
				return d.getTime() / 1000;
			},
		}`,
	}
}
