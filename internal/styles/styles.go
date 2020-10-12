/*
Copyright © 2020 DSCVIT

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package styles

import "github.com/ttacon/chalk"

// WarnStyle style for decorating WARN text
var WarnStyle = chalk.Bold.NewStyle().WithBackground(chalk.Yellow).WithForeground(chalk.Black)

// InfoStyle style for decorating INFO text
var InfoStyle = chalk.Bold.NewStyle().WithBackground(chalk.ResetColor).WithForeground(chalk.Green)

// ErrorStyle style for decorating ERROR text
var ErrorStyle = chalk.Bold.NewStyle().WithBackground(chalk.ResetColor).WithForeground(chalk.Red)
