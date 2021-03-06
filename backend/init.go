/*
 * @Copyright Reserved By Janusec (https://www.janusec.com/).
 * @Author: U2
 * @Date: 2018-07-14 16:22:46
 * @Last Modified: U2, 2018-07-14 16:22:46
 */

package backend

import (
	"janusec/data"
	"janusec/utils"

	_ "github.com/lib/pq"
)

func InitDatabase() {
	dal := data.DAL
	err := dal.CreateTableIfNotExistsCertificates()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsApplications()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsDomains()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsDestinations()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsSettings()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsAppUsers()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	_, err = dal.InsertIfNotExistsAppUser(`admin`, `1f7d7e9decee9561f457bbc64dd76173ea3e1c6f13f0f55dc1bc4e99e5b8b494`,
		`afa8bae009c9dbf4135f62e165847227`, ``, true, true, true, true)
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsNodes()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	err = dal.CreateTableIfNotExistsTOTP()
	if err != nil {
		utils.DebugPrintln("InitDatabase", err)
	}
	// Upgrade to latest version
	if dal.ExistColumnInTable("domains", "redirect") == false {
		// v0.9.6+ required
		err = dal.ExecSQL(`alter table domains add column redirect boolean default false, add column location varchar(256) default null`)
		if err != nil {
			utils.DebugPrintln("InitDatabase", err)
		}
	}
	if dal.ExistColumnInTable("applications", "oauth_required") == false {
		// v0.9.7+ required
		err = dal.ExecSQL(`alter table applications add column oauth_required boolean default false, add column session_seconds bigint default 7200, add column owner varchar(128)`)
		if err != nil {
			utils.DebugPrintln("InitDatabase", err)
		}
	}
	if dal.ExistColumnInTable("destinations", "route_type") == false {
		// v0.9.8+ required
		err = dal.ExecSQL(`alter table destinations add column route_type bigint default 1, add column request_route varchar(128) default '/', add column backend_route varchar(128) default '/'`)
		if err != nil {
			utils.DebugPrintln("InitDatabase", err)
		}
	}
	if dal.ExistColumnInTable("ccpolicies", "interval_seconds") == true {
		// v0.9.9 interval_seconds, v0.9.10 interval_milliseconds
		err = dal.ExecSQL(`ALTER TABLE ccpolicies RENAME COLUMN interval_seconds TO interval_milliseconds`)
		if err != nil {
			utils.DebugPrintln("InitDatabase", err)
		}
		err = dal.ExecSQL(`UPDATE ccpolicies SET interval_milliseconds=interval_milliseconds*1000`)
		if err != nil {
			utils.DebugPrintln("InitDatabase", err)
		}
	}
}

func LoadAppConfiguration() {
	LoadCerts()
	LoadApps()
	if data.IsPrimary {
		LoadDestinations()
		LoadDomains()
		LoadAppDomainNames()
		LoadNodes()
	} else {
		LoadRoute()
		LoadDomains()
	}
}
